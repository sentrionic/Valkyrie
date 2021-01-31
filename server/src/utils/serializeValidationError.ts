import { ValidationError } from 'yup';

export type ValidationErrors = {
  field: string;
  message: string;
};

export const serializeValidationError = (
  err: ValidationError,
): ValidationErrors[] => {
  const invalid: ValidationErrors[] = [];

  err.inner.map((value) => {
    invalid.push({
      field: value.path!,
      message: value.errors[0],
    });
  });
  return invalid;
};