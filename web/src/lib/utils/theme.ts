import { extendTheme } from "@chakra-ui/react";
import { mode } from "@chakra-ui/theme-tools";

const config: any = {
  initialColorMode: 'dark',
};

const styles = {
  global: (props: any) => ({
    body: {
      bg: mode("gray.100", "#1b1c1d")(props),
    },
  }),
};

const colors = {
  highlight: {
    standard: "#7289da",
    hover: "#677bc4",
    active: "#5b6eae",
  },
  brandGray: {
    accent: "#8e9297",
    light: "#36393f",
    dark: "#303339",
  },
};

const fonts = {
  body: "'Open Sans', sans-serif"
}

const customTheme = extendTheme({ colors, config, styles, fonts });

export default customTheme;

export const scrollbarCss = {
  "&::-webkit-scrollbar": {
    width: "8px",
  },
  "&::-webkit-scrollbar-track": {
    background: "#2f3136",
    width: "10px",
  },
  "&::-webkit-scrollbar-thumb": {
    background: "#202225",
    borderRadius: "18px",
  },
}
