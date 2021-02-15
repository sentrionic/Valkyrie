import { InternalServerErrorException } from '@nestjs/common';
import * as aws from 'aws-sdk';
import * as path from 'path';
import * as sharp from 'sharp';
import { config } from 'dotenv';
import { BufferFile } from '../types/BufferFile';
import { nanoid } from 'nanoid';

config();

/**
 * Shrinks the image to 150 x 150 and turns it into .webp
 * @param buffer
 */
const imageTransformer = (buffer: Buffer): Promise<Buffer> =>
  sharp(buffer)
    .resize({
      width: 150,
      height: 150,
    })
    .webp()
    .toBuffer();

const s3 = new aws.S3({
  accessKeyId: process.env.AWS_ACCESS_KEY,
  secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
  region: process.env.AWS_S3_REGION,
});

/**
 * Resizes an image and uploads it to S3
 * @param directory
 * @param image
 */
export const uploadAvatarToS3 = async (
  directory: string,
  image: BufferFile,
): Promise<string> => {
  const { buffer } = await image;
  const stream = await imageTransformer(buffer);

  if (!stream) {
    throw new InternalServerErrorException();
  }

  const params = {
    Bucket: process.env.AWS_STORAGE_BUCKET_NAME as string,
    Key: `files/${directory}/${nanoid(20)}.webp`,
    Body: stream,
    ContentType: 'image/webp',
  };

  const response = await s3.upload(params).promise();

  return response.Location;
};

/**
 * Uploads a file to S3
 * @param directory
 * @param file
 */
export const uploadToS3 = async (
  directory: string,
  file: BufferFile,
): Promise<string> => {
  const { buffer, originalname, mimetype } = await file;

  const params = {
    Bucket: process.env.AWS_STORAGE_BUCKET_NAME as string,
    Key: `files/${directory}/${formatName(originalname)}`,
    Body: buffer,
    ContentType: mimetype,
  };

  const response = await s3.upload(params).promise();

  return response.Location;
};

const formatName = (filename: string): string => {
  const file = path.parse(filename);
  const name = file.name;
  const ext = file.ext
  const date = Date.now();
  const cleanFileName = name.toLowerCase().replace(/[^a-z0-9]/g, '-');
  return `${date}-${cleanFileName}${ext}`;
};

export const deleteFile = async (filename: string): Promise<void> => {
  const index = filename.indexOf('files');
  const key = filename.slice(index);

  const params = {
    Bucket: process.env.AWS_STORAGE_BUCKET_NAME as string,
    Key: key,
  };

  s3.deleteObject(params, (err) => {
    if (err) console.log(err, err.stack);
  });
};
