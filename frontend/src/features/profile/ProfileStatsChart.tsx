import React from 'react';
import { PieChart, Pie, ResponsiveContainer, Tooltip } from 'recharts';
import { renderCustomTooltip } from './ProfileToolTip';
import type { StatItem } from './profileTypes';
import styles from './Profile.module.css';

interface ProfileStatsChartProps {
    data: StatItem[];
}

interface PieData extends StatItem {
    fill: string;
}

const mapToPieData = (item: StatItem): PieData => {
    return {
        name: item.name,
        value: item.value,
        color: item.color,
        fill: item.color
    };
};

const transformDataForPie = (data: StatItem[]): PieData[] => {
    return data.map(mapToPieData);
};

export const ProfileStatsChart: React.FC<ProfileStatsChartProps> = ({ data }) => {
    return (
        <div className={styles.chartSection}>
            <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                    <Tooltip content={renderCustomTooltip} />
                    <Pie
                        data={transformDataForPie(data)}
                        dataKey="value"
                        nameKey="name"
                        cx="50%"
                        cy="50%"
                        innerRadius={40}
                        outerRadius={65}
                        paddingAngle={4}
                        stroke="none"
                    />
                </PieChart>
            </ResponsiveContainer>
        </div>
    );
};