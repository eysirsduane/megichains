import dayjs, { locale } from 'dayjs';
import 'dayjs/locale/zh-cn';
import 'dayjs/locale/en';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import { localStg } from '@/utils/storage';

/**
 * Set dayjs locale
 *
 * @param lang
 */
export function setDayjsLocale(lang: App.I18n.LangType = 'zh-CN') {
  const localMap = {
    'zh-CN': 'zh-cn',
    'en-US': 'en'
  } satisfies Record<App.I18n.LangType, string>;

  const l = lang || localStg.get('lang') || 'zh-CN';

  locale(localMap[l]);
}

export function getHumannessDateTime(unix: number) {
  if (unix === 0) {
    return '';
  }

  dayjs.extend(utc);
  dayjs.extend(timezone);

  return dayjs.tz(unix, 'Asia/Yerevan').format('YYYY-MM-DD HH:mm:ss');
}
