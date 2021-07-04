import {Client4} from 'mattermost-redux/client';
import {ClientError} from 'mattermost-redux/client/client4';

import manifest from 'manifest';

const apiUrl = `/plugins/${manifest.id}/api/v0`;

export const setChannelAutoFollow = async (channelId: string, autoFollow: boolean) => {
    console.log("Set autofollow for channel " + channelId + " to " + autoFollow);

    const method = autoFollow ? 'put' : 'delete';

    // await doFetchWithResponse(apiUrl + '/autofollow/' + channelId, {method: method});
}

export const getChannelAutoFollow = async (channelId: string) => {
    // const {data} = await doFetchWithResponse(apiUrl + '/autofollow/' + channelId, {method: 'get'});

    return true;
}

export const doFetchWithResponse = async (url: string, options = {}) => {
    const response = await fetch(url, Client4.getOptions(options));

    let data;
    if (response.ok) {
        data = await response.json();
        return {
            response,
            data,
        };
    }

    data = await response.text();

    throw new ClientError(Client4.url, {
        message: data || '',
        status_code: response.status,
        url,
    });
};
