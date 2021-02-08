import * as yup from 'yup';

export const ChannelSchema = yup.object().shape({
  name: yup.string().min(3).max(30).required(),
  isPublic: yup.boolean().optional().default(true),
  members: yup
    .array(
      yup.string().optional().max(20, 'Must provide memberId'),
    )
    .optional()
});
