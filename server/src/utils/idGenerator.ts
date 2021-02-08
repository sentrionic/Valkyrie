import { customAlphabet } from 'nanoid';

const alphabet = '0123456789';
const generator = customAlphabet(alphabet, 20);

export const idGenerator = (): string => generator();
