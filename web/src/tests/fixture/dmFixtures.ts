import { DMChannel } from '../../lib/models/dm';

export const mockDMChannel: DMChannel = {
  id: '1446384585456750592',
  user: {
    id: '1446384528997224448',
    username: 'Alice',
    image: 'https://gravatar.com/avatar/c160f8cc69a4f0bf2b0362752353d060?d=identicon',
    isOnline: false,
    isFriend: false,
  },
};

export const mockDMChannelList = [mockDMChannel];
