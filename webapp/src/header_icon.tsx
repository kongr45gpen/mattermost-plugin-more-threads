import React, {FC, useState, useEffect} from 'react';
import styled from 'styled-components';
import {useSelector, useDispatch} from 'react-redux';

import {GlobalState} from 'mattermost-redux/types/store';
import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';

import {getChannelAutoFollow} from 'client';

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
    const [isAutoFollowing, setAutoFollowing] = useState(false);

    useEffect(() => {
        const doFetch = async () => {
            const autoFollowing = await getChannelAutoFollow(currentChannelId);
            setAutoFollowing(autoFollowing);
        };
        doFetch();
    }, [currentChannelId, dispatch]);

    if (isAutoFollowing) {
        return (
            <div>
                <HeaderIconContainer
                    active={true}
                    className={'fa fa-heart'}
                />
                Autofollowing
            </div>
        );
    }

    return (
        <HeaderIconContainer
            className={'fa fa-heart-o'}
        />
    );
};

export default HeaderIcon;
