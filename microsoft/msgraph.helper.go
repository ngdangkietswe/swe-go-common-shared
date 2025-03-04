package microsoft

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/samber/lo"
)

type MSGraphHelper struct {
	credential interface{}
	appClient  *msgraphsdk.GraphServiceClient
}

func NewMSGraphHelper() *MSGraphHelper {
	return &MSGraphHelper{}
}

// InitializeMSGraph initializes the MSGraphHelper for application authentication
func (g *MSGraphHelper) InitializeMSGraph(isDeviceCodeCred bool) error {
	var (
		credential interface{}
		scopes     []string
		err        error
	)

	// Build the credential based on the authentication type (client secret or device code)
	if isDeviceCodeCred {
		credential, err = g.buildDeviceCodeCredential()
		scopes = []string{
			"https://graph.microsoft.com/Chat.Create",
			"https://graph.microsoft.com/ChatMessage.Send",
			"https://graph.microsoft.com/User.Read",
		}
	} else {
		credential, err = g.buildClientSecretCredential()
	}

	if err != nil {
		return err
	}

	g.credential = credential

	// Create an auth provider with the client secret credential
	authProvider, err := g.buildAuthProvider(scopes)
	if err != nil {
		return err
	}

	// Create a request adapter with the auth provider
	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return err
	}

	// Create a Graph client with the request adapter
	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.appClient = client

	return nil
}

// buildClientSecretCredential builds a client secret credential
func (g *MSGraphHelper) buildClientSecretCredential() (*azidentity.ClientSecretCredential, error) {
	return azidentity.NewClientSecretCredential(
		config.GetString("MAAD_TENANT_ID", "s3cur3d"),
		config.GetString("MAAD_CLIENT_ID", "s3cur3d"),
		config.GetString("MAAD_CLIENT_SECRET", "s3cur3d"),
		nil,
	)
}

// buildDeviceCodeCredential builds a device code credential
func (g *MSGraphHelper) buildDeviceCodeCredential() (*azidentity.DeviceCodeCredential, error) {
	userPrompt := func(ctx context.Context, options azidentity.DeviceCodeMessage) error {
		fmt.Println(options.Message)
		return nil
	}

	return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		TenantID:   config.GetString("MAAD_TENANT_ID", "s3cur3d"),
		ClientID:   config.GetString("MAAD_CLIENT_ID", "s3cur3d"),
		UserPrompt: userPrompt,
	})
}

// buildAuthProvider builds an auth provider
func (g *MSGraphHelper) buildAuthProvider(scopes []string) (*auth.AzureIdentityAuthenticationProvider, error) {
	var cred azcore.TokenCredential

	switch g.credential.(type) {
	case *azidentity.ClientSecretCredential:
		cred = g.credential.(*azidentity.ClientSecretCredential)
	case *azidentity.DeviceCodeCredential:
		cred = g.credential.(*azidentity.DeviceCodeCredential)
	default:
		return nil, fmt.Errorf("unsupported credential type")
	}

	return auth.NewAzureIdentityAuthenticationProviderWithScopes(cred, scopes)
}

// GetUserIDByEmail gets the user ID by email
func (g *MSGraphHelper) GetUserIDByEmail(email string) (string, error) {
	filter := fmt.Sprintf("mail eq '%s'", email)
	queryParams := &users.UsersRequestBuilderGetQueryParameters{
		Filter: &filter,
		Select: []string{"id"},
	}

	cfg := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: queryParams,
	}

	result, err := g.appClient.Users().Get(context.Background(), cfg)
	if err != nil {
		return "", err
	}

	msTeamUsers := result.GetValue()
	if len(msTeamUsers) == 0 {
		return "", fmt.Errorf("user not found with email: %s", email)
	}

	userID := msTeamUsers[0].GetId()
	if userID == nil {
		return "", fmt.Errorf("user not found with email: %s", email)
	}

	return *userID, nil
}

// initMember initializes a member with the given ID
func (g *MSGraphHelper) initMember(memberEmail string) *models.AadUserConversationMember {
	member := models.NewAadUserConversationMember()
	member.SetRoles([]string{"owner"})
	member.SetAdditionalData(map[string]interface{}{
		"user@odata.bind": "https://graph.microsoft.com/v1.0/users/" + memberEmail,
	})

	return member
}

// initChat initializes a chat with the given sender and receiver IDs
func (g *MSGraphHelper) initChat(senderEmail, receiverEmail string) *models.Chat {
	chat := models.NewChat()
	chat.SetChatType(lo.ToPtr(models.ONEONONE_CHATTYPE))
	chat.SetMembers(
		[]models.ConversationMemberable{
			g.initMember(senderEmail),
			g.initMember(receiverEmail),
		})

	return chat
}

// SendMessage sends a message from the sender to the receiver with the given content
func (g *MSGraphHelper) SendMessage(ctx context.Context, senderEmail, receiverEmail, content string) error {
	createdChat, err := g.appClient.Chats().Post(ctx, g.initChat(senderEmail, receiverEmail), nil)
	if err != nil {
		return err
	}

	chatMsg := models.NewChatMessage()

	body := models.NewItemBody()
	body.SetContent(lo.ToPtr(content))
	body.SetContentType(lo.ToPtr(models.HTML_BODYTYPE))

	chatMsg.SetBody(body)

	_, err = g.appClient.Chats().ByChatId(*createdChat.GetId()).Messages().Post(ctx, chatMsg, nil)
	if err != nil {
		return err
	}

	return nil
}
