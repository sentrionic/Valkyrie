import { ValidationError } from 'yup';

export type ValidationErrors = {
  field: string;
  message: string;
};

/**
 * Creates an error array of the format
 * {
 *   field: field name [e.g. username],
 *   message: error message [e.g. "The username is too short"]
 * }
 * @param err
 */
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
