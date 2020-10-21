# Cloudformation

## 1. Sample command

To create this stack using the `template.yml` file attached:

1. Copy the template file to s3:

(The command below is run from the project root directory)

```sh
aws s3 cp cloudformation/template.yml s3://${BUCKET_NAME}/template-$(date +%F).yml
```

2. (Optional) Once the template file has been uploaded to s3, validate the template:

```sh
aws cloudformation validate-template --template-url https://s3.me-south-1.amazonaws.com/${BUCKET_NAME}/template-2020-10-18.yml
```

3. Create the stack passing in the required parameters:

```sh
aws cloudformation create-stack --stack-name short-url \
        --template-url https://s3.me-south-1.amazonaws.com/${BUCKET_NAME}/template-2020-10-18.yml \
        --parameters \
        ParameterKey=VPC,ParameterValue=${VPC_ID} \
        ParameterKey=PublicSubnetA,ParameterValue=${PUBLIC_SUBNET_A_ID} \
        ParameterKey=PublicSubnetB,ParameterValue=${PUBLIC_SUBNET_B_ID} \
        ParameterKey=PrivateSubnetA,ParameterValue=${PRIVATE_SUBNET_A_ID} \
        ParameterKey=PrivateSubnetB,ParameterValue=${PRIVATE_SUBNET_B_ID} \
        ParameterKey=IamCertificateArn,ParameterValue=${IAM_CERTIFICATE_ARN} \
        ParameterKey=IamCertificateId,ParameterValue=${IAM_CERTIFICATE_SERVER_CERTIFICATE_ID}  \
        ParameterKey=Image,ParameterValue=${IMAGE_URI} \
        ParameterKey=ServiceName,ParameterValue=short-url \
        ParameterKey=HealthCheckPath,ParameterValue=/health \
        ParameterKey=HostedZoneName,ParameterValue=kurtz.ml \
        ParameterKey=DatabaseConnectionString,ParameterValue=${DATABASE_CONNECTION_STRING} \
        ParameterKey=ContainerSecurityGroup,ParameterValue=sg-${CONTAINER_SECURITY_GROUP} \
        ParameterKey=LoadBalancerSecurityGroup,ParameterValue=${LOAD_BALANCER_SECURITY_GROUP} \
        ParameterKey=BaseUrl,ParameterValue=https://kurtz.ml  \
        --capabilities CAPABILITY_NAMED_IAM \
        --region me-south-1
```

## 2. Parameters

**IAM Certificate**

- CloudFront only supports ACM certificates in the US East (N. Virginia) Region (us-east-1). [Source](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudfront-distribution-viewercertificate.html#cfn-cloudfront-distribution-viewercertificate-acmcertificatearn)
- In unsupported Regions, you must use IAM as a certificate manager
- [This](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_server-certs.html) guide documents how to upload an IAM Certificate.

- You can obtain an SSL Certificate from ZeroSSL or Let's Encrypt.
- The certificate contains a `certificate.crt`, `ca.bundle.crt` and `private.key` files.
- IAM Certificates require the certificates to be in PEM format.
- To convert the `.crt` files to `.pem`, run the following command:

  ```sh
  openssl x509 -in certificate.crt -out certificate.pem -outform PEM
  ```

  ```sh
  openssl x509 -in ca_bundle.crt -out ca_bundle.pem -outform PEM
  ```

- To convert the `.key` file to `.pem`, run the following command:

  ```
  openssl rsa -in private.key -out private.pem -outform PEM
  ```

- Then run the following command to create an IAM certificate:

  ```
  aws iam upload-server-certificate --server-certificate-name ExampleCertificateName --certificate-body file:////certificate.pem --certificate-chain file:///ca_bundle.pem --private-key file:///private.pem --path /cloudfront/
  ```

  **NOTE**: We use `file://` because the command expects the contents of the file and not paths to the file.

- The returned has the following structure:

  ```json
  {
    "ServerCertificateMetadata": {
      "Path": "/cloudfront/",
      "ServerCertificateName": "ExampleCertificateName",
      "ServerCertificateId": "ASCA4GIXYZM4VYZIJKLMN",
      "Arn": "arn:aws:iam::12345678:server-certificate/cloudfront/ExampleCertificateName",
      "UploadDate": "2020-10-16T13:00:33+00:00",
      "Expiration": "2021-01-14T23:59:59+00:00"
    }
  }
  ```

- The `Arn` field's value is provided to the LoadBalancer and the `ServerCertificateId` field's value is provided to Cloudfront.

## 3. IAM Permissions

The minimal set of IAM permissions to run the template file are specified in the policy below:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "iam:GetRole",
        "iam:DetachRolePolicy",
        "iam:CreateRole",
        "iam:DeleteRole",
        "iam:AttachRolePolicy",
        "logs:PutRetentionPolicy"
      ],
      "Resource": [
        "arn:aws:iam::<AccountId>:role/*",
        "arn:aws:logs:*:<AccountId>:log-group:*"
      ]
    }
  ]
}
```

## 4. Notes & Future Improvements

1. The Fargate containers run on the public subnet (meaning the parameters `PrivateSubnetA` and `PrivateSubnetB` are unused). This is because the Fargate contaienrs need to download the image and therefore, if they were in a private subnet, a NAT gateway would be needed for the network connectivity. Rather than pay for a NAT Gateway, I'm running the Fargate containers in the public subnet.

2. A classical load balancer would be cheaper than the application load balancer used here.
