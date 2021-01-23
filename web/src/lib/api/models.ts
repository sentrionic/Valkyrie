export interface FieldError {
  field: string;
  message: string;
}

export interface Member {
  id: string;
  username: string;
  avatar: string;
  createdAt: string;
  updatedAt: string;
}

export interface Message {
  id: string;
  text?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AccountResponse {
  id: string;
  username: string;
  email: string;
  image: string;
  createdAt: string;
  updatedAt: string;
}
