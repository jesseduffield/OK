import * as aws from '@pulumi/aws';
import * as pulumi from '@pulumi/pulumi';
import * as mime from 'mime';
import * as fs from 'fs';
import * as path from 'path';

// prerequisite: you must have okquestionmark.org as a hosted zone in route53. The www subdomain will be created by pulumi
const domain = 'www.okquestionmark.org';
const hostedZoneId = 'Z042093996TWBRN7V284';

const siteBucket = new aws.s3.Bucket(domain, {
  website: {
    indexDocument: 'index.html',
  },
  // need to set this so that pulumi doesn't append it's own SHA suffix: the
  // website needs to share the name of the S3 bucket.
  bucket: domain,
});

const siteDir = '../../site/build'; // directory for content files

// from https://gist.github.com/lovasoa/8691344
const walkSync = (dir: string, callback: (filepath: string) => void) => {
  const files = fs.readdirSync(dir);
  files.forEach(file => {
    var filepath = path.join(dir, file);
    const stats = fs.statSync(filepath);
    if (stats.isDirectory()) {
      walkSync(filepath, callback);
    } else if (stats.isFile()) {
      callback(filepath);
    }
  });
};

walkSync(siteDir, filepath => {
  const relativePath = path.relative(siteDir, filepath);
  new aws.s3.BucketObject(relativePath, {
    bucket: siteBucket,
    source: new pulumi.asset.FileAsset(filepath), // use FileAsset to point to a file
    contentType: mime.getType(filepath) || undefined, // set the MIME type of the file
  });
});

exports.bucketName = siteBucket.bucket; // create a stack export for bucket name

// Create an S3 Bucket Policy to allow public read of all objects in bucket
// This reusable function can be pulled out into its own module
function publicReadPolicyForBucket(bucketName: string) {
  return JSON.stringify({
    Version: '2012-10-17',
    Statement: [
      {
        Effect: 'Allow',
        Principal: '*',
        Action: ['s3:GetObject'],
        Resource: [
          `arn:aws:s3:::${bucketName}/*`, // policy refers to bucket name explicitly
        ],
      },
    ],
  });
}

// Set the access policy for the bucket so all objects are readable
const bucketPolicy = new aws.s3.BucketPolicy('bucketPolicy', {
  bucket: siteBucket.bucket, // depends on siteBucket -- see explanation below
  policy: siteBucket.bucket.apply(publicReadPolicyForBucket),
  // transform the siteBucket.bucket output property -- see explanation below
});

// add record
const record = new aws.route53.Record(domain, {
  name: domain,
  zoneId: hostedZoneId,
  type: 'CNAME',
  ttl: 60,
  records: [siteBucket.websiteEndpoint],
});

exports.websiteUrl = siteBucket.websiteEndpoint; // output the endpoint as a stack output
