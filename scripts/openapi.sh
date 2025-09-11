openapi-generator-cli generate \
  -i openapi/cmc-rest.yaml \
  -g go-server \
  -o /tmp/oapi \
  --additional-properties=packageName=cmcrest

rm -rf openapi/servergen
mkdir openapi/servergen
mv /tmp/oapi/go/ openapi/servergen
rm -rf /tmp/oapi

openapi-generator-cli generate \
  -i openapi/cmc-rest.yaml \
  -g go \
  -o /tmp/oapi \
  --additional-properties=packageName=cmcrest

rm -rf openapi/clientgen
mkdir -p openapi/clientgen/go
mv /tmp/oapi/*.go openapi/clientgen/go
rm -rf /tmp/oapi
