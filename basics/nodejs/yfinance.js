const yahooFinance = require('yahoo-finance2').default; // NOTE the .default

async function getStockData(symbol) {
  try {
    const results = await yahooFinance.quote(symbol);
    const json = JSON.stringify(results);
    console.log(json);
  } catch (error) {
    console.log(error)
    console.error("Error fetching data:", error.message);
    process.exit(1);
  }
}

const symbol = process.argv[2];
if (!symbol) {
  console.error("Please provide a Stock symbol as an argument");
  process.exit(1);
}

getStockData(symbol);
