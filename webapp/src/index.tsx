import {Store, Action} from 'redux';

import {GlobalState} from 'mattermost-redux/types/store';

import HeaderIcon from 'header_icon';
import {setChannelAutoFollow} from 'client';

// eslint-disable-next-line import/no-unresolved
import {PluginRegistry} from './types/mattermost-webapp';

import manifest from './manifest';

export default class Plugin {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars, @typescript-eslint/no-empty-function
    public async initialize(registry: PluginRegistry, store: Store<GlobalState, Action<Record<string, unknown>>>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

		// eslint-disable-next-line TS2686 TS7006
        registry.registerChannelHeaderButtonAction(<HeaderIcon/>, (channel) => {
            setChannelAutoFollow(channel.id, true);
        }, 'Auto-follow threads', 'Auto-follow threads');
    }
}

declare global {
    interface Window {
        registerPlugin(id: string, plugin: Plugin): void
    }
}

window.registerPlugin(manifest.id, new Plugin());
