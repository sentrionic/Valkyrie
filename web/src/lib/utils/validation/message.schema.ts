import * as yup from 'yup';

const SUPPORTED_FORMATS = ['image/jpg', 'image/jpeg', 'audio/mp3', 'audio/mpeg', 'image/png'];

export const FileSchema = yup.object().shape({
  file: yup
    .mixed()
    .nullable()
    .test('fileSize', 'The file is too large', (value) => value?.size < 5000000)
    .test(
      'type',
      'Only the following formats are accepted: Image and Audio',
      (value) => value && SUPPORTED_FORMATS.includes(value.type)
    ),
});
