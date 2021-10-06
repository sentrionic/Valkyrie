import dayjs from 'dayjs';
import calender from 'dayjs/plugin/calendar';

dayjs.extend(calender);

export const getTime = (createdAt: string): string => dayjs(createdAt).calendar();

export const getShortenedTime = (createdAt: string): string => dayjs(createdAt).format('h:mm A');

export const getTimeDifference = (date1: string, date2: string): number => dayjs(date1).diff(dayjs(date2), 'minutes');

export const checkNewDay = (date1: string, date2: string): boolean => !dayjs(date1).isSame(dayjs(date2), 'day');

export const formatDivider = (date: string): string => dayjs(date).format('MMMM D, YYYY');
