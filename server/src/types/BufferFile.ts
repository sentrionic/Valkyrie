/**
 * Represents an uploaded file
 */
export type BufferFile = {
  fieldname: string;
  originalname: string;
  encoding: string;
  mimetype: string;
  buffer: Buffer;
  size: number;
};
