#!/bin/bash
echo "Starting altron..."

function ctrl_c() {
    echo "Stopping altron..."
    docker compose down
    exit 1
}

trap ctrl_c INT

echo "Enter altron scale number:"
read scale_number

export ALTRON_REPLICAS=$scale_number

# create ftp server with n users
container_id=$(docker run -dt \
    -p 21:21 \
    -p 30000-30009:30000-30009 \
    -e FTP_USER_HOME=/home/enmex \
    -e FTP_USER_NAME=enmex \
    -e FTP_USER_PASS=9T9m0jt*InSy \
    -e FTP_MAX_CLIENTS=15 \
    stilliard/pure-ftpd:latest)

for i in $(seq 1 15) 
do
    user=$(echo "altron$i")
    docker exec -it $container_id \
        pure-pw useradd $user \
            -f /etc/pure-ftpd/passwd/pureftpd.passwd \
            -m \
            -u ftpuser \
            -d /home/$user
done

# starting app
#sudo docker compose up --build -d frontend
kubectl apply -f manifest.yml
kubectl apply -f service.yml
