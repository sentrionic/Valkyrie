import dayjs from 'dayjs';
import calender from 'dayjs/plugin/calendar';

dayjs.extend(calender);

export const getTime = (createdAt: string): string => {
  return dayjs(createdAt).calendar();
};