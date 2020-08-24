#!/bin/sh

ENV_VARS="BUCKET_NAME BUCKET_LIFECYCLE_DAYS WEB_CLIENT_S3_ACCESS_KEY WEB_CLIENT_S3_SECRET_KEY GRPC_SERVER_S3_ACCESS_KEY GRPC_SERVER_S3_SECRET_KEY"

# --- input validation ---
for var in $ENV_VARS;
do
	if [ -z "$(eval echo \$$var)" ]
	then
		echo "Error: $var is not defined!"
		return 1
	fi
done

echo -n "Setting up '$BUCKET_NAME' bucket and users... "

# --- bucket configuration ---

# evaluate lifecycle days within bucket policy
cd $(dirname "$0")
envsubst < data/.minio.sys/buckets/BUCKET_NAME/lifecycle.xml > /tmp/lifecycle.xml
mv -f /tmp/lifecycle.xml data/.minio.sys/buckets/BUCKET_NAME/lifecycle.xml

# change placeholder bucket name to actual bucket name, create bucket folder
mv data/.minio.sys/buckets/BUCKET_NAME data/.minio.sys/buckets/$BUCKET_NAME
mkdir data/$BUCKET_NAME

# --- user configuration ---

# setup new users
cd data/.minio.sys/config/iam
for user in WEB_CLIENT GRPC_SERVER;
do
	# evaluate credentials within identity files
	envsubst < users/$user/identity.json > /tmp/identity.json
	mv -f /tmp/identity.json users/$user/identity.json

	# change placeholder user directory names to actual user IDs
	mv users/$user users/$(eval "echo \$${user}_S3_ACCESS_KEY")

	# change placeholder policy file names to actual user IDs 
	mv policydb/users/$user.json policydb/users/$(eval "echo \$${user}_S3_ACCESS_KEY").json
done

# evaluate bucket name within canned policies for users
for policy in images-portal-grpc-server-policy images-portal-web-client-policy;
do
	envsubst < policies/$policy/policy.json > /tmp/policy.json
	mv -f /tmp/policy.json policies/$policy/policy.json
done
cd - > /dev/null

# --- apply configuration ---
cp -rf data/. /data

echo "Done!"

echo "Launching minio..."
/usr/bin/docker-entrypoint.sh server /data
