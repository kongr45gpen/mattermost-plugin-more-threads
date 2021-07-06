import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import HeaderIcon from 'header_icon';
import {setChannelAutoFollow} from 'client';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

import manifest from './manifest';

import {statuses} from './statuses';

export default class Plugin {

    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

		let uniqid = 0;

		const handler = (channel: any) => {
			let [isAutofollowing, setAutofollow] = statuses[channel.id];

			const newFollowStatus = !isAutofollowing;

			setChannelAutoFollow(channel.id, newFollowStatus).then(() => {
				uniqid = registry.registerChannelHeaderButtonAction(<HeaderIcon/>, handler, 'Auto-follow threads', 'Auto-follow threads');
			});

			// Terrible hacks to set UI autofollow
			setAutofollow(newFollowStatus)
			statuses[channel.id][0] = newFollowStatus;
			registry.unregisterComponent(uniqid);
			console.log("Setting auto-follow status to " + newFollowStatus);
        }

		// eslint-disable-next-line
        uniqid = registry.registerChannelHeaderButtonAction(<HeaderIcon/>, handler, 'Auto-follow threads', 'Auto-follow threads');
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
