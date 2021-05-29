import {
  ArgumentMetadata,
  HttpException,
  HttpStatus,
  Injectable,
  PipeTransform,
} from '@nestjs/common';
import { serializeValidationError } from './serializeValidationError';

@Injectable()
export class YupValidationPipe implements PipeTransform {
  constructor(private readonly schema: any) {}

  async transform(value: any, metadata: ArgumentMetadata): Promise<any> {
    try {
      await this.schema.validate(value, { abortEarly: false });
    } catch (err) {
      const errors = serializeValidationError(err);
      throw new HttpException(
        { message: 'Input data validation failed', errors },
        HttpStatus.BAD_REQUEST,
      );
    }
    return value;
  }
}
