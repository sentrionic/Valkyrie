import * as yup from 'yup';

export const MessageSchema = yup.object().shape({
  text: yup.string().optional().test("empty", "Message must not be empty", text => text?.length !== 0),
});