from aws_cdk import aws_s3 as s3
from aws_cdk import aws_ssm as ssm
from aws_cdk import core


class InfraStack(core.Stack):
    def __init__(self, scope: core.Construct, construct_id: str, **kwargs) -> None:
        super().__init__(scope, construct_id, **kwargs)

        base_bucket_name = "svl-joke-bot-bucket"
        bucket = s3.Bucket(self, base_bucket_name)

        # persist bucket.bucket_name to Param Store with key `/Stage/{base_bucket_name}`
        ssm.StringParameter(
            self,
            "Parameter",
            description="S3 deploy bucket for svl-joke-bot",
            parameter_name=f"/Stage/{base_bucket_name}",
            string_value=bucket.bucket_name,
        )
