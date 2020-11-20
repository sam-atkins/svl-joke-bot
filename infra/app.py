#!/usr/bin/env python3
import os

from aws_cdk import core

from infra.infra_stack import InfraStack

REGION = os.getenv("CDK_DEFAULT_REGION", "eu-west-2")

app = core.App()
InfraStack(app, "infra", env={"region": REGION})

app.synth()
