import React, { InputHTMLAttributes } from "react";
import { useField } from "formik";
import {
  FormControl,
  FormErrorMessage,
  FormLabel,
  Text,
  Input,
} from "@chakra-ui/react";

type InputFieldProps = InputHTMLAttributes<HTMLInputElement> & {
  label: string;
  name: string;
};

export const InputField: React.FC<InputFieldProps> = ({ label, ...props }) => {
  const [field, { error, touched }] = useField(props);
  return (
    <FormControl mt={4} isInvalid={error != null && touched}>
      <FormLabel htmlFor={field.name}>
        <Text fontSize="12px" textTransform="uppercase">
          {label}
        </Text>
      </FormLabel>
      {/* @ts-ignore */}
      <Input
        bg="brandGray.dark"
        borderColor="black"
        borderRadius="3px"
        focusBorderColor="highlight.standard"
        {...field}
        {...props}
      />
      <FormErrorMessage>{error}</FormErrorMessage>
    </FormControl>
  );
};
