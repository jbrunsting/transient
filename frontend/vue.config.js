module.exports = {
    devServer: {
        host: '0.0.0.0',
        watchOptions: {
            ignored: /node_modules/,
            aggregateTimeout: 300,
            poll: 1000,
        },
    },
};
