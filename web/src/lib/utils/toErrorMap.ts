import { FieldError } from '../api/models';

export const toErrorMap = (errors: FieldError[]) => {
  const errorMap: Record<string, string> = {};
  errors.forEach(({ field, message }) => {
    errorMap[field.toLowerCase()] = message;
  });

  return errorMap;
};
