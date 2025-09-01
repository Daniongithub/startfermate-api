const express = require('express');
const axios = require('axios');
const cheerio = require('cheerio');
const cors = require('cors');

const app = express();
const port = 3005;
const version = "2.1.6";

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
            linea = linea.replace("Linea ", "");
            var destinazione = element.find('.bus-destination').text().trim();
            var mezzo = element.find('.det a').attr('data-vehicle') || '';

            //Aggiustamenti

            if(mezzo.length == 4){
              mezzo = 3 + mezzo;
            }

            if(stato == "Non disp"){
              stato = "Non disponibile";
            }

            if(mezzo == ""){
              mezzo = "Non disponibile";
            }

            if(destinazione == "Fornace.Zarattini"){
              destinazione = "Fornace Zarattini"
            }

            //Varianti linee

            if(linea == "1" && destinazione == "Borgo Nuovo"){
              linea = "1B";
            }

            if(linea == "4" && destinazione == "Mirabilandia"){
              linea = "4B";
            }

            if(linea == "4" && destinazione == "Classe Cantoniera"){
              linea = "4C";
            }

            if(linea == "4" && destinazione == "ClasseCantoniera"){
              linea = "4C";
            }

            if(linea == "4" && destinazione == "Lido di Dante"){
              linea = "4D";
            }

            if(linea == "4" && destinazione == "Fosso Ghiaia"){
              linea = "4F";
            }

            if(linea == "4" && destinazione == "Classe via Liburna"){
              linea = "4R";
            }

            if(linea == "4" && destinazione == "Classe Romea Vecchia"){
              linea = "4R";
            }

            //Linee limitate o soppresse a metà

            if(linea == "8" && destinazione == "Deposito"){
              linea = "8/";
            }

            if(linea == "3" && destinazione == "via Sant'Alberto"){
              linea = "3/";
            }

            const linee = ["1", "1B", "3", "4", "4B", "4D", "5", "8", "18", "70", "80"];

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

app.get('/versione', async (req, res) => {
  res.send(version);
});

app.listen(port, () => {
    console.log(`API attiva su http://localhost:${port}`);
});
