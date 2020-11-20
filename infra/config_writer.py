import os
import sys

import boto3
from dotenv import load_dotenv
from tomlkit import dumps, parse

load_dotenv()


REGION = os.getenv("CDK_DEFAULT_REGION")
CDK_PROFILE_NAME = os.getenv("AWS_CDK_PROFILE")
SVL_PROFILE_NAME = os.getenv("AWS_SVL_PROFILE")

DOC = parse(
    """version = 0.1

[default.build.parameters]
profile = ""
debug = true
region = ""
skip_pull_image = true
use_container = false

[default.local_start_api.parameters]
port = 8080

[default.package.parameters]
profile = "sam__dev_svl"
region = "eu-west-2"
output_template_file = "packaged.yaml"
s3_bucket = ""

[default.deploy.parameters]
profile = ""
capabilities = "CAPABILITY_IAM"
region = ""
s3_bucket = ""
"""
)


def _session():
    session = boto3.Session(profile_name=CDK_PROFILE_NAME, region_name=REGION)
    return session.client("ssm")


def get_param(param_name: str) -> str:
    session = _session()
    try:
        response = session.get_parameter(Name=param_name)
    except Exception as ex:
        print(f"Exception with boto request: {ex}")
        sys.exit(1)

    try:
        status_code = response.get("ResponseMetadata").get("HTTPStatusCode")
    except AttributeError:
        print("Malformed response from boto request to Param Store failed")
        sys.exit(1)

    if status_code != 200:
        print("Something went wrong with the request to the Param Store")
        sys.exit(1)

    try:
        parameter = response.get("Parameter").get("Value")
    except AttributeError:
        print("Unable to get Param value from Param Store response")
        sys.exit(1)

    return parameter


# TODO(sam) pass in arg e.g. --bucket-name which is provided to write_toml_config
def write_toml_config():
    bucket_name = get_param(param_name="/Stage/svl-joke-bot-bucket")

    # set build details
    DOC["default"]["build"]["parameters"]["profile"] = SVL_PROFILE_NAME
    DOC["default"]["build"]["parameters"]["region"] = REGION

    # set package details
    DOC["default"]["package"]["parameters"]["profile"] = SVL_PROFILE_NAME
    DOC["default"]["package"]["parameters"]["region"] = REGION
    DOC["default"]["package"]["parameters"]["s3_bucket"] = bucket_name

    # set deploy details
    DOC["default"]["deploy"]["parameters"]["profile"] = SVL_PROFILE_NAME
    DOC["default"]["deploy"]["parameters"]["region"] = REGION
    DOC["default"]["deploy"]["parameters"]["s3_bucket"] = bucket_name

    with open("../samconfig.toml", "w") as f:
        f.write(dumps(DOC))

    print("✅ samconfig.toml file written")


write_toml_config()
