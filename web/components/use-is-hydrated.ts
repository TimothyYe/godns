'use client';

import { useSyncExternalStore } from 'react';

const subscribe = () => () => {};

export const useIsHydrated = () => {
	return useSyncExternalStore(subscribe, () => true, () => false);
};
