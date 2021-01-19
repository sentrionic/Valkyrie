import { InternalServerErrorException } from '@nestjs/common';
import * as aws from 'aws-sdk';
import * as path from 'path';
import * as sharp from 'sharp';
import { config } from 'dotenv';
import { BufferFile } from '../types/BufferFile';

config();

const DIM_MAX = 1080;
const DIM_MIN = 320;

const storyImageTransformer = async (
  buffer: Buffer,
): Promise<Buffer | null> => {
  const image = sharp(buffer);
  return await new Promise<Buffer | null>(async (res) =>
    image.metadata().then((info) => {
      res(image.resize(getResizeOptions(info)).webp().toBuffer());
    }),
  );
};

const getResizeOptions = (info: sharp.Metadata): sharp.ResizeOptions => {
  let options: Record<string, unknown>;
  if (
    (info.height !== undefined && info.height < DIM_MIN) ||
    (info.width !== undefined && info.width < DIM_MIN)
  ) {
    options = {
      width: DIM_MIN,
      height: DIM_MIN,
      fit: 'outside',
    };
  } else {
    options = {
      width: DIM_MAX,
      height: DIM_MAX,
      fit: 'inside',
      withoutEnlargement: true,
    };
  }

  return options;
};

const profileImageTransformer = (buffer: Buffer): Promise<Buffer> =>
  sharp(buffer)
    .resize({
      width: 150,
      height: 150,
    })
    .webp({
      quality: 75,
    })
    .toBuffer();

const s3 = new aws.S3({
  accessKeyId: process.env.AWS_ACCESS_KEY,
  secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
  region: process.env.AWS_S3_REGION,
});

export const uploadAvatarToS3 = async (
  directory: string,
  image: BufferFile,
): Promise<string> => {
  const { buffer } = await image;
  const stream = await profileImageTransformer(buffer);

  if (!stream) {
    throw new InternalServerErrorException();
  }

  const params = {
    Bucket: process.env.AWS_STORAGE_BUCKET_NAME as string,
    Key: `files/${directory}/avatar.webp`,
    Body: stream,
    ContentType: 'image/webp',
  };

  const response = await s3.upload(params).promise();

  return response.Location;
};

export const uploadToS3 = async (
  directory: string,
  image: BufferFile,
): Promise<string> => {
  const { buffer, originalname } = await image;

  const stream = await storyImageTransformer(buffer);

  if (!stream) {
    throw new InternalServerErrorException();
  }

  const params = {
    Bucket: process.env.AWS_STORAGE_BUCKET_NAME as string,
    Key: `files/${directory}/${formatName(originalname)}`,
    Body: stream,
    ContentType: 'image/webp',
  };

  const response = await s3.upload(params).promise();

  return response.Location;
};

const formatName = (filename: string): string => {
  const name = path.parse(filename).name;
  const date = Date.now();
  const cleanFileName = name.toLowerCase().replace(/[^a-z0-9]/g, '-');
  return `${date}-${cleanFileName}.webp`;
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