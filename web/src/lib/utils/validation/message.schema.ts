import * as yup from 'yup';

const SUPPORTED_FORMATS = ['image/jpg', 'image/jpeg', 'audio/mp3', 'image/png'];

export const FileSchema = yup.object().shape({
  file: yup
    .mixed()
    .nullable()
    .test('fileSize', 'The file is too large', (value) => {
      return value?.size < 5000000;
    })
    .test('type', 'Only the following formats are accepted: Image and Audio', (value) => {
      return value && SUPPORTED_FORMATS.includes(value.type);
    }),
});
