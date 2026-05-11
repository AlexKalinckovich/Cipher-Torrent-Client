import React from 'react';
import { Table, Typography, Alert } from 'antd';
import type {ColumnsType} from 'antd/es/table';
import type {Torrent} from '../../types/model/models.ts';
import { useTorrents } from '../../hooks/useTorrents.ts';

export const Dashboard: React.FC = () => {
    console.log("Use torrents")
    const { data, isLoading, error } = useTorrents();

    console.log(data, isLoading, error)
    return renderContent(data ?? [], isLoading, error);
};

const renderContent = (torrents: Torrent[], isLoading: boolean, error: Error | null): React.ReactElement => {
    if (error) {
        return renderError(error);
    }
    return renderTable(torrents, isLoading);
};

const renderError = (error: Error): React.ReactElement => {
    return <Alert title={error.message} type="error" showIcon />;
};

const renderTable = (torrents: Torrent[], isLoading: boolean): React.ReactElement => {
    return (
        <div style={{ padding: '24px' }}>
            <Typography.Title level={2}>Dashboard</Typography.Title>
            <Table<Torrent>
                dataSource={torrents}
                columns={getColumns()}
                rowKey="info_hash"
                loading={isLoading}
                pagination={false}
            />
        </div>
    );
};

const getColumns = (): ColumnsType<Torrent> => {
    return [
        { title: 'Name', dataIndex: 'name', key: 'name' },
        { title: 'Status', dataIndex: 'status', key: 'status' },
        { title: 'Progress (%)', dataIndex: 'progress', key: 'progress' },
        { title: 'Size (Bytes)', dataIndex: 'size_bytes', key: 'size_bytes' },
        { title: 'Peers', dataIndex: 'peers_count', key: 'peers_count' },
    ];
};