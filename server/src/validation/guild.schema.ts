import * as yup from 'yup';

export const GuildSchema = yup.object().shape({
  name: yup.string().min(3).max(30).required(),
});

export const UpdateGuildSchema = yup.object().shape({
  name: yup.string().nullable().min(3).max(30),
});
