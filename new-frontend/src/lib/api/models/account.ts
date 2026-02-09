import type { AccountNoID } from './accountNoID';
import type { Entity } from './index';

export type Account = Entity & AccountNoID & Required<Pick<Entity & AccountNoID, 'showInDashboardSummary'>>;
