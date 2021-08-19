import * as aws from '@pulumi/aws';
import * as pulumi from '@pulumi/pulumi';
import * as awsx from '@pulumi/awsx';

// see https://www.pulumi.com/docs/reference/pkg/aws/lambda/function/
// and https://github.com/pulumi/examples/blob/master/aws-go-lambda

const lambdaName = 'ok-run';

const role = new aws.iam.Role('task-exec-role', {
  assumeRolePolicy: JSON.stringify({
    Version: '2012-10-17',
    Statement: [
      {
        Action: 'sts:AssumeRole',
        Effect: 'Allow',
        Sid: '',
        Principal: {
          Service: 'lambda.amazonaws.com',
        },
      },
    ],
  }),
});

const logPolicy = new aws.iam.Policy('lambda-log-policy', {
  policy: JSON.stringify({
    Version: '2012-10-17',
    Statement: [
      {
        Action: [
          'logs:CreateLogGroup',
          'logs:CreateLogStream',
          'logs:PutLogEvents',
        ],
        Resource: 'arn:aws:logs:*:*:*',
        Effect: 'Allow',
      },
    ],
  }),
});

const rolePolicyAttachment = new aws.iam.RolePolicyAttachment('lambda-logs', {
  role: role,
  policyArn: logPolicy.arn,
});

let fileArchive = new pulumi.asset.FileArchive('./handler/handler.zip');

const func = new aws.lambda.Function(
  lambdaName,
  {
    handler: 'handler',
    role: role.arn,
    runtime: 'go1.x',
    code: fileArchive,
  },
  {
    dependsOn: [rolePolicyAttachment],
  }
);

// Define a new GET endpoint from an existing Lambda Function.
const api = new awsx.apigateway.API('ok-api-prod', {
  routes: [
    {
      path: '/run',
      method: 'POST',
      eventHandler: func,
    },
    {
      // see https://github.com/pulumi/pulumi-awsx/issues/545#issuecomment-880579206
      path: '/run',
      method: 'OPTIONS',
      eventHandler: async () => {
        return {
          body: '',
          statusCode: 200,
          headers: {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Credentials': 'true',
            'Access-Control-Allow-Methods': 'POST, OPTIONS',
            'Access-Control-Allow-Headers':
              'Origin, X-Requested-With, Content-Type, Accept, Authorization',
          },
        };
      },
    },
  ],
});

// Export the auto-generated API Gateway base URL.
export const url = api.url;
