import * as yup from 'yup';

export const MemberSchema = yup.object().shape({
  nickname: yup
    .string()
    .nullable()
    .min(3)
    .max(30),
  color: yup
    .string()
    .nullable()
    .matches(/^#[0-9a-f]{3}(?:[0-9a-f]{3})?$/i, "The color must be a valid hex color")
});