import { customAlphabet } from 'nanoid';

const alphabet = '0123456789';
const generator = customAlphabet(alphabet, 20);

/**
 * Generates a 20 character long numeric snowflake id
 */
export const idGenerator = (): string => generator();
