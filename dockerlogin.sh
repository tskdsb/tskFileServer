IPS=$(kubectl get no -o go-template='{{range .items}}{{with index .status.addresses 0}}{{.address}} {{end}}{{end}}')
for ip in ${IPS}; do
  #./tskssh_darwin_amd64 -ip ${ip} -password caicloud2019 -cmd "docker login cargo-infra.caicloud.xyz -u admin -p C2njlZ1vsgI8"
  ./tskssh_darwin_amd64 -ip ${ip} -password caicloud2019 -cmd "cp /root/.docker/config.json /var/lib/kubelet/"
done
