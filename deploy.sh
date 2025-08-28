echo "Enter build tag: "
read build_tag
echo "Building Altron version '$build_tag'"

docker compose build

docker tag altron-core enmex/altron:core-$build_tag
docker tag altron-session enmex/altron:session-$build_tag
docker tag altron-plugin enmex/altron:plugin-$build_tag
docker tag altron-frontend enmex/altron:frontend-$build_tag
docker tag altron-converter enmex/altron:converter-$build_tag
docker tag altron-connection enmex/altron:connection-$build_tag

docker push enmex/altron:core-$build_tag
docker push enmex/altron:session-$build_tag
docker push enmex/altron:plugin-$build_tag
docker push enmex/altron:frontend-$build_tag
docker push enmex/altron:converter-$build_tag
docker push enmex/altron:connection-$build_tag