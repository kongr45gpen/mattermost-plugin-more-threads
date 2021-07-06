import React, {FC, useState, useEffect} from 'react';
import styled from 'styled-components';
import {useSelector, useDispatch} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {getChannelAutoFollow} from 'client';

import {statuses} from 'statuses';

interface HeaderIconProps {
    active?: boolean;
}

const HeaderIconContainer = styled.div<HeaderIconProps>`
    font-size: 18px;
    ${(props) => (props.active ? 'color: var(--button-bg);' : '')}
`;

const HeaderIcon:FC = () => {
    const dispatch = useDispatch();
    const currentChannelId = useSelector<GlobalState, string>(getCurrentChannelId);
    let [isAutoFollowing, setAutoFollowing] = useState(false);

    const headericon = this;

    useEffect(() => {
        const doFetch = async () => {
            console.log("Asked to fetch autofollow status");

            const autoFollowing = await getChannelAutoFollow(currentChannelId);
            setAutoFollowing(autoFollowing);
            statuses[currentChannelId] = [autoFollowing, setAutoFollowing];
        };
        doFetch();
    }, [currentChannelId, dispatch]);

    if (isAutoFollowing) {
        return (
            <span
                className={'Button Button___transparent is-active separated FollowButton'}
            >
                <span
                    className='Button_label'
                >
                    {'Auto-following'}
                </span>
            </span>
        );
    }

    return (
        <span
            className={'Button Button___transparent separated FollowButton'}
        >
            <span
                className='Button_label'
            >
                {'Auto-follow'}
            </span>
        </span>
    );
};

export default HeaderIcon;
