echo "Update protobuf..."
git clone https://github.com/ngdangkietswe/swe-protobuf-shared.git
cp -r swe-protobuf-shared/generated/common proto/
rm -rf swe-protobuf-shared
echo "Update protobuf successful!"