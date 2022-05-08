#!/bin/sh

CERT_REQUEST_CONFIG="certRequest.config"
EXT_FILE="v3.ext"

generationType=
option=
caKeyPem=
caCertPem=
caDuration=

certDir=
certKeyPem=
certRequestFile=
certPem=
certDuration=

#FQDN of the server
CN=
#Country Name (2 letter code)
C=
#Locality Name (eg, city)
L=
#Organization Name (eg, company)
O=
#Organizational Unit Name (eg, section)
OU= 
SAN=

usage() {
    echo "Generate CA"
    echo "Example: "
    echo ""
    echo "./tls.sh generate --generationType CA  --certDir ./cfg/tls/sec --caKeyPem rootCAkey.pem --caCertPem rootCACert.pem --caDuration 3650"
    echo ""
    echo ""
    echo "Generate Certificate signed by the CA"
    echo "Example: "
    echo ""
    echo "./tls.sh generate --generationType cert --certDir ./cfg/tls/sec --caKeyPem rootCAkey.pem --caCertPem rootCACert.pem --certDuration 365 --CN localhost --C BG --L Sofia --O TT --OU system --SAN DNS:localhost,DNS:sensorapi --certRequestFile serverCert.csr --certPem serverCert.pem --certKeyPem serverKey.pem"
    echo ""
    echo ""
}

while [ "$1" != "" ]; do
    case $1 in
        generate | deploy ) option=$1
            ;;
        --generationType )
            shift
            generationType=$1
            ;;
        --caKeyPem )
            shift
            caKeyPem=$1
            ;;
        --caCertPem )
            shift
            caCertPem=$1
            ;;
        --caDuration )
            shift
            caDuration=$1
            ;;
         --certDir )
            shift
            certDir=$1
            ;;
        --certKeyPem )
            shift
            certKeyPem=$1
            ;;
        --certRequestFile )
            shift
            certRequestFile=$1
            ;;
        --certPem )
            shift
            certPem=$1
            ;;
          --certDuration )
            shift
            certDuration=$1
            ;;
        --CN )
            shift
            CN=$1
            ;;
        --C )
            shift
            C=$1
            ;;
        --L )
            shift
            L=$1
            ;;
        --O )
            shift
            O=$1
            ;;
        --OU )
            shift
            OU=$1
            ;;
        --SAN )
            shift
            SAN=$1
            ;;
        -h | --help )
            usage
            exit
            ;;
        * )
            usage
            exit 1
    esac
    shift
done

generateRootCA() {
    mkdir -p $certDir
    openssl genrsa -out "$certDir/$caKeyPem" 2048
    openssl req -x509 -sha256 -new -nodes -key "$certDir/$caKeyPem" -days "$caDuration" -out "$certDir/$caCertPem"
}

generateCertificate() {
    printf '[req]
req_extensions = v3_req
distinguished_name = dn
prompt = no

[dn]
CN = %s
C = %s
L = %s
O = %s
OU = %s

[v3_req]
subjectAltName = %s' "$CN" "$C" "$L" "$O" "$OU" "$SAN" >> "$certDir/$CERT_REQUEST_CONFIG"

#Create certificate private key with the request
openssl genrsa -out "$certDir/$certKeyPem" 2048
openssl req -new -key "$certDir/$certKeyPem" -sha256 -out "$certDir/$certRequestFile" -config "$certDir/$CERT_REQUEST_CONFIG"

printf 'subjectAltName = %s' "$SAN" >> "$certDir/$EXT_FILE"
openssl x509 -req -sha256 -in "$certDir/$certRequestFile" -CA "$certDir/$caCertPem" -CAkey "$certDir/$caKeyPem" -CAcreateserial -out "$certDir/$certPem" -days "$certDuration" -extfile "$certDir/$EXT_FILE"

rm -r "$certDir/$CERT_REQUEST_CONFIG" "$certDir/$EXT_FILE"
}

if [ "$option" == 'generate' ] && [ "$generationType" == 'CA' ]
then
    generateRootCA
fi

if [ "$option" == 'generate' ] && [ "$generationType" == 'cert' ]
then
    generateCertificate
fi