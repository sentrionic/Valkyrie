import * as yup from 'yup';

export const ChannelSchema = yup.object().shape({
  name: yup.string().min(3).max(30).required(),
  isPublic: yup.boolean().default(true),
  // members: yup
  //   .array()
  //   .optional()
  //   .min(1, 'At least one member')
  //   .defined(),
});