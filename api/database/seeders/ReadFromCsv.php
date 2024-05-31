<?php

namespace Database\Seeders;

class ReadFromCsv
{
    public static function getDataFromCsv($path, $separator = null, $enclosure = '"')
    {
        $file = fopen($path,'r');
        $header = fgetcsv($file, separator: $separator, enclosure: $enclosure);
        $data = [];

        while ($row = fgetcsv($file, separator: $separator, enclosure: $enclosure)) {
            $data[] = array_combine($header, $row);
        }

        fclose($file);
        return $data;
    }
}
