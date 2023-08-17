version=$1
docker build -t  jialiannexus.xgd.com:8084/k8s-addon/alert-rlist:$version .

docker push jialiannexus.xgd.com:8084/k8s-addon/alert-rlist:$version
