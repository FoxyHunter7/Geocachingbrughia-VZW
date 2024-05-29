<?php

namespace Database\Seeders;

use App\Models\MessageLanguage;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class MessageLanguageSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/messagesLanguage.csv'), ';');

        $model = new MessageLanguage();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
