const express = require('express');
const axios = require('axios');
const cheerio = require('cheerio');
const cors = require('cors');

const app = express();
const port = 3005;

app.use(cors());
app.use(express.json());

app.get('/fermata', async (req, res) => {
    try {
        const { param, param2, palina } = req.query;
        const url = `https://infobus.startromagna.it/InfoFermata?param=${param}&param2=${param2}&palina=${palina}`;
        const response = await axios.get(url);
        const $ = cheerio.load(response.data);
        const results = [];
        var fermata = $('h2.fw-bold.text-primary.title').text().trim();
        if(fermata != ""){
          results.push({
          fermata
          });
        }
        $('.container.mb-50 .bus-card').each((i, el) => {
            const element = $(el);
            const isSoppressa = element.find('.bus-status.sopp').length > 0;
            const headerSpan = element.find('.bus-header span').first();
            const orario = element.find('.bus-times span').first().text().trim();
            var stato = element.find('.bus-status').text().trim();
            var linea = headerSpan.contents().filter((i, el) => el.type === 'text').text().trim();
            const destinazione = element.find('.bus-destination').text().trim();
            var mezzo = element.find('.det a').attr('data-vehicle') || '';

            if(mezzo.length == 4){
              mezzo = 3 + mezzo;
            }

            if(stato == "Non disp"){
              stato = "Non disponibile";
            }

            if(mezzo == ""){
              mezzo = "Non disponibile";
            }

            //Varianti linee

            if(linea == "Linea 4" && destinazione == "Mirabilandia"){
              linea = "Linea 4B";
            }

            if(linea == "Linea 4" && destinazione == "Lido di Dante"){
              linea = "Linea 4D";
            }

            if(linea == "Linea 4" && destinazione == "Classe via Liburna"){
              linea = "Linea 4R";
            }

            if(linea == "Linea 4" && destinazione == "Classe Romea Vecchia"){
              linea = "Linea 4R";
            }

            if(linea == "Linea 1" && destinazione == "Borgo Nuovo"){
              linea = "Linea 1B";
            }

            //Linee soppresse a metà

            const linee = ["Linea 1", "Linea 1B", "Linea 3", "Linea 4", "Linea 4B", "Linea 4D", "Linea 5", "Linea 8", "Linea 18", "Linea 70", "Linea 80"];

            if(linee.includes(linea) && destinazione == "Stazione FS"){
              linea = linea + "/";
            }

            results.push({
                linea,
                destinazione,
                orario,
                stato,
                mezzo,
                soppressa: isSoppressa
            });
        });

        res.json(results);
    } catch (error) {
        console.error('Errore:', error.message);
        res.status(500).send('Errore nel recupero dei dati');
    }
});

app.get('/bacino', async (req, res) => {
    const selectedOption = req.query.selectedOption;

    if (!selectedOption || !['ra', 'rn', 'fc'].includes(selectedOption)) {
        return res.status(400).json({ error: 'Parametro selectedOption mancante o non valido' });
    }

    try {
        const response = await axios.post('https://infobus.startromagna.it/FermateService.asmx/GetFermateByBacino', {
            bacino: selectedOption
        }, {
            headers: {
                'Content-Type': 'application/json'
            }
        });

        const raw = response.data.d;
        const data = raw.map(({ nome, palina, targetID }) => ({ nome, palina, targetID }));
        if (Array.isArray(data)) {
          res.json(data);
          } else {
          res.status(500).json({ error: 'Formato di risposta non valido dal server remoto' });
          }
    } catch (error) {
        console.error('Errore nella richiesta al servizio remoto:', error);
        res.status(500).json({ error: 'Errore nel contattare il servizio remoto' });
    }
});
app.listen(port, () => {
    console.log(`API attiva su http://localhost:${port}`);
});
