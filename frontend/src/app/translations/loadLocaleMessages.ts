import type { AppLocale } from './detectLocale';

import ptCommon from './pt-BR/common.json';
import ptLogin from './pt-BR/login.json';
import ptCode from './pt-BR/code.json';
import ptNav from './pt-BR/nav.json';
import ptHome from './pt-BR/home.json';
import ptTokenSetup from './pt-BR/tokenSetup.json';
import ptApiKey from './pt-BR/apiKey.json';
import ptTokenReveal from './pt-BR/tokenReveal.json';
import ptUsers from './pt-BR/users.json';
import ptAudit from './pt-BR/audit.json';
import ptTeams from './pt-BR/teams.json';
import ptTags from './pt-BR/tags.json';
import ptEnvironments from './pt-BR/environments.json';
import ptProjects from './pt-BR/projects.json';
import ptSecrets from './pt-BR/secrets.json';

import enCommon from './en-US/common.json';
import enLogin from './en-US/login.json';
import enCode from './en-US/code.json';
import enNav from './en-US/nav.json';
import enHome from './en-US/home.json';
import enTokenSetup from './en-US/tokenSetup.json';
import enApiKey from './en-US/apiKey.json';
import enTokenReveal from './en-US/tokenReveal.json';
import enUsers from './en-US/users.json';
import enAudit from './en-US/audit.json';
import enTeams from './en-US/teams.json';
import enTags from './en-US/tags.json';
import enEnvironments from './en-US/environments.json';
import enProjects from './en-US/projects.json';
import enSecrets from './en-US/secrets.json';

const ptBR = {
  common: ptCommon,
  login: ptLogin,
  code: ptCode,
  nav: ptNav,
  home: ptHome,
  tokenSetup: ptTokenSetup,
  apiKey: ptApiKey,
  tokenReveal: ptTokenReveal,
  users: ptUsers,
  audit: ptAudit,
  teams: ptTeams,
  tags: ptTags,
  environments: ptEnvironments,
  projects: ptProjects,
  secrets: ptSecrets,
};

const enUS = {
  common: enCommon,
  login: enLogin,
  code: enCode,
  nav: enNav,
  home: enHome,
  tokenSetup: enTokenSetup,
  apiKey: enApiKey,
  tokenReveal: enTokenReveal,
  users: enUsers,
  audit: enAudit,
  teams: enTeams,
  tags: enTags,
  environments: enEnvironments,
  projects: enProjects,
  secrets: enSecrets,
};

export type MessageSchema = typeof ptBR;

export function loadLocaleMessages(locale: AppLocale): MessageSchema {
  switch (locale) {
    case 'pt-BR':
      return ptBR;
    case 'en-US':
    default:
      return enUS;
  }
}

export function allLocaleMessages(): Record<AppLocale, MessageSchema> {
  return {
    'pt-BR': ptBR,
    'en-US': enUS,
  };
}
