#!/usr/bin/env python3
import os

from aws_cdk import core

from dotenv import load_dotenv
from infra.infra_stack import InfraStack

load_dotenv()

REGION = os.getenv("CDK_DEFAULT_REGION", "eu-west-2")

app = core.App()
InfraStack(app, "infra", env={"region": REGION})

app.synth()
