import * as yup from 'yup';

const SUPPORTED_FORMATS = ['image/jpg', 'image/jpeg', 'audio/mp3', 'audio/mpeg', 'image/png'];

export const FileSchema = yup.object().shape({
  file: yup
    .mixed<FileList>()
    .nullable()
    .test('count', 'Only one file is allowed', (value) => value?.length === 1)
    .test('fileSize', 'The file is too large', (value) => !!value && value[0].size < 5000000)
    .test(
      'type',
      'Only the following formats are accepted: Image and Audio',
      (value) => !!value && SUPPORTED_FORMATS.includes(value[0].type)
    ),
});
